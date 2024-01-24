package cred

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/credm"
	"github.com/techierishi/pal/logr"
	palSync "github.com/techierishi/pal/sync"
	"github.com/techierishi/pal/util"
	"github.com/techierishi/pal/wrapper"
)

var credNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new credential",
	RunE:  newFunc,
}

func init() {
	RootCmdCred.AddCommand(credNewCmd)
}

func newFunc(cmd *cobra.Command, args []string) (err error) {
	logger := logr.GetLogInstance()

	credDb, err := credm.GetCredDb(logger)
	defer credDb.Close()
	if err != nil {
		return err
	}
	credentials := credm.Credentials{
		Db:     credDb,
		Logger: logger,
	}
	credentialMap, err := credentials.Load()
	if err != nil {
		return err
	}

	newCredential := credm.NewCred()

	hash := util.CalculateHash(fmt.Sprintf("%s_%s", newCredential.Application, newCredential.Username))

	if _, ok := credentialMap[hash]; ok {
		return fmt.Errorf("Credential for [%s] already exists", newCredential.Application)
	}

	keyRing := wrapper.KeyRing{Logger: logger}
	ok := keyRing.Set(newCredential.Application, newCredential.Username, newCredential.Password)
	if !ok {
		fmt.Printf("%s\n", color.RedString("Keyring not supported on this platform!"))
	}

	// Password is only saved in os keychain
	newCredential.Password = ""
	newCredential.Timestamp = util.UnixMilli()
	credDb.Upsert(fmt.Sprintf("%s.%s", credm.CRED_TBL, hash), newCredential)
	fmt.Fprintf(color.Output, "%12s", color.GreenString("Credential saved! \n"))

	if config.Conf.Gist.AutoSync {
		return palSync.AutoSync(config.Conf.General.BackupFile)
	}

	return nil
}
