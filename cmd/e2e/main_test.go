package e2e

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/caarlos0/env/v6"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/anoriar/gophkeeper/cmd/e2e/dto/response"
)

const (
	SuccessMessage = "success"
	FailMessage    = "fail"
)

type GophkeeperSuite struct {
	suite.Suite

	serverProcess *exec.Cmd
	conf          *Config

	clientStdErr *os.File
}

func (suite *GophkeeperSuite) SetupSuite() {
	conf := NewConfig()

	err := env.Parse(conf)

	suite.NoError(err)
	suite.Require().NotEmpty(conf.ServerBin)
	suite.Require().NotEmpty(conf.ServerDatabaseURI)
	suite.Require().NotEmpty(conf.ServerRunAddress)
	suite.Require().NotEmpty(conf.ServerPublicAddress)
	suite.Require().NotEmpty(conf.ClientBin)
	suite.Require().NotEmpty(conf.ClientDataDirName)

	suite.conf = conf

	suite.startServerProcess(*conf)
	suite.prepareClient()
}

func (suite *GophkeeperSuite) prepareClient() {
	//настройка переменных окружения
	err := os.Setenv("SERVER_ADDRESS", suite.conf.ServerPublicAddress)
	require.NoError(suite.T(), err)
	err = os.Setenv("DATA_DIRNAME", suite.conf.ClientDataDirName)
	require.NoError(suite.T(), err)

	//создание директории для хранения данных
	err = os.MkdirAll(suite.conf.ClientDataDirName, 0755)
	assert.NoError(suite.T(), err, "Error creating directory")

	//создание логов для клиента
	clientStderr, err := os.Create("./stderr-client.log")
	assert.NoError(suite.T(), err)
	suite.clientStdErr = clientStderr

}

func (suite *GophkeeperSuite) startServerProcess(conf Config) {
	serverSystemEnvs := append(os.Environ(),
		"RUN_ADDRESS="+conf.ServerRunAddress,
		"DATABASE_URI="+conf.ServerDatabaseURI,
		"JWT_SECRET_KEY="+"secret",
	)

	suite.serverProcess = exec.Command(conf.ServerBin)
	suite.serverProcess.Env = append(suite.serverProcess.Env, serverSystemEnvs...)

	stdout, err := os.Create("./stdout-server.log")
	assert.NoError(suite.T(), err)

	defer stdout.Close()
	suite.serverProcess.Stdout = stdout

	stderr, err := os.Create("./stderr-server.log")
	assert.NoError(suite.T(), err)

	defer stderr.Close()
	suite.serverProcess.Stderr = stderr

	err = suite.serverProcess.Start()
	assert.NoError(suite.T(), err, "Error starting the binary")

	time.Sleep(2 * time.Second)
}

func (suite *GophkeeperSuite) TearDownSuite() {
	time.Sleep(2 * time.Second)

	if suite.serverProcess != nil && suite.serverProcess.Process != nil {
		_ = suite.serverProcess.Process.Kill()
	}

	db, err := sqlx.Open("pgx", suite.conf.ServerDatabaseURI)
	assert.NoError(suite.T(), err, "Error tear down database")

	err = goose.DownTo(db.DB, "internal/server/migrations", 0)
	assert.NoError(suite.T(), err, "Error tear down database")

	err = os.RemoveAll(suite.conf.ClientDataDirName)
	assert.NoError(suite.T(), err, "Error removing directory")

	suite.clientStdErr.Close()
}

func (suite *GophkeeperSuite) TestGophkeeper() {
	defer suite.TearDownSuite()

	entriesMap := make(map[string]response.LoginDetailEntryResponse)

	suite.Run("register user", func() {
		commandResponse := suite.execCommand("register", "-u", "test22", "-p", "1234567", "-m", "2345")
		assert.Equal(suite.T(), SuccessMessage, commandResponse.Status)
	})
	suite.Run("login user", func() {

		commandResponse := suite.execCommand("login", "-u", "test22", "-p", "1234567", "-m", "2345")
		assert.Equal(suite.T(), SuccessMessage, commandResponse.Status)
	})

	addedId := ""
	suite.Run("add login entry", func() {
		commandResponse := suite.execCommand("add", "-t", "login", "-d", "{\"login\": \"test\", \"password\": \"pass\"}", "-m", "{\"prop1\": \"val1\", \"prop2\": \"val2\"}")
		assert.Equal(suite.T(), SuccessMessage, commandResponse.Status)
		var loginPayload response.LoginDetailEntryResponse
		err := json.Unmarshal(commandResponse.Payload, &loginPayload)
		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), "test", loginPayload.Data.Login)
		assert.Equal(suite.T(), "pass", loginPayload.Data.Password)

		addedId = loginPayload.Id
		entriesMap[loginPayload.Id] = loginPayload
	})
	suite.Run("edit login entry", func() {
		commandResponse := suite.execCommand("edit", "-t", "login", "-i", addedId, "-d", "{\"login\": \"test2\", \"password\": \"pass2\"}", "-m", "{\"prop1\": \"val1\", \"prop2\": \"val2\"}")
		assert.Equal(suite.T(), SuccessMessage, commandResponse.Status)
		var loginPayload response.LoginDetailEntryResponse
		err := json.Unmarshal(commandResponse.Payload, &loginPayload)
		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), "test2", loginPayload.Data.Login)
		assert.Equal(suite.T(), "pass2", loginPayload.Data.Password)

		entriesMap[loginPayload.Id] = loginPayload
	})

	suite.Run("sync login entry", func() {
		commandResponse := suite.execCommand("sync", "-t", "login")
		assert.Equal(suite.T(), SuccessMessage, commandResponse.Status)
	})

	suite.Run("detail login entry", func() {
		commandResponse := suite.execCommand("detail", "-t", "login", "-i", addedId)
		assert.Equal(suite.T(), SuccessMessage, commandResponse.Status)
		var loginPayload response.LoginDetailEntryResponse
		err := json.Unmarshal(commandResponse.Payload, &loginPayload)
		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), "test2", loginPayload.Data.Login)
		assert.Equal(suite.T(), "pass2", loginPayload.Data.Password)
	})
	suite.Run("list login entry", func() {
		commandResponse := suite.execCommand("list", "-t", "login")
		assert.Equal(suite.T(), SuccessMessage, commandResponse.Status)
		var loginListPayload []response.LoginListEntryResponse
		err := json.Unmarshal(commandResponse.Payload, &loginListPayload)
		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), 1, len(loginListPayload))
		assert.Equal(suite.T(), addedId, loginListPayload[0].Id)
	})

	suite.Run("delete login entry", func() {
		commandResponse := suite.execCommand("delete", "-t", "login", "-i", addedId)
		assert.Equal(suite.T(), SuccessMessage, commandResponse.Status)
	})

	suite.Run("sync login entry", func() {
		commandResponse := suite.execCommand("sync", "-t", "login")
		assert.Equal(suite.T(), SuccessMessage, commandResponse.Status)
	})

	suite.Run("detail login entry", func() {
		commandResponse := suite.execCommand("detail", "-t", "login", "-i", addedId)
		assert.Equal(suite.T(), FailMessage, commandResponse.Status)
	})

	suite.Run("list login entry", func() {
		commandResponse := suite.execCommand("list", "-t", "login")
		assert.Equal(suite.T(), SuccessMessage, commandResponse.Status)
		var loginListPayload []response.LoginListEntryResponse
		err := json.Unmarshal(commandResponse.Payload, &loginListPayload)
		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), 0, len(loginListPayload))
	})
}

func (suite *GophkeeperSuite) execCommand(arg ...string) response.CommandResponse {
	var stdout bytes.Buffer
	cmd := suite.prepareCommand(&stdout, arg...)
	err := cmd.Run()
	require.NoError(suite.T(), err)
	var commandResponse response.CommandResponse
	err = json.Unmarshal(stdout.Bytes(), &commandResponse)
	require.NoError(suite.T(), err)
	return commandResponse
}

func (suite *GophkeeperSuite) prepareCommand(stdout *bytes.Buffer, arg ...string) *exec.Cmd {
	cmd := exec.Command(suite.conf.ClientBin, arg...)
	cmd.Env = os.Environ()
	cmd.Stdout = stdout
	cmd.Stderr = suite.clientStdErr
	return cmd
}

func TestMyTestSuite(t *testing.T) {
	_, exists := os.LookupEnv("CLIENT_BIN")
	if !exists {
		t.Skip()
	}
	suite.Run(t, new(GophkeeperSuite))
}
