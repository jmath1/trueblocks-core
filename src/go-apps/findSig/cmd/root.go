package cmd

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func scrapeBlocks() {
	var search [][]byte
	for _, i := range os.Args[1:] {
		decoded, err := hex.DecodeString(i[2:])
		if err != nil {
			fmt.Printf("Couldn't parse argument %s\n", i)
		} else {
			search = append(search, decoded)
		}
	}

	// Create thread pool with number of concurrent threads equal to 'runtime.NumCPU()'
	var wait sync.WaitGroup
	checkOne, _ := ants.NewPoolWithFunc(runtime.NumCPU(), func(canonical interface{}) {
		// Calculate 4-byte form
		res := crypto.Keccak256([]byte(canonical.(string)))[:4]
		// Go through queries and compare
		for i := 0; i < len(search); i++ {
			cur := search[i]
			if bytes.Equal(res, cur) {
				fmt.Printf("                                                            \r")
				fmt.Printf("\"0x%s\",\"%s\"\n", hex.EncodeToString(res), canonical)
				os.Exit(0)
			}
			fmt.Fprintf(os.Stderr, "\033[2KScanning: %s\r", canonical)
		}
		wait.Done()
	})

	// Load files
	fTemp2, _ := os.Open("/Users/jrush/.quickBlocks/known_abis/uniq_funcs.tab")
	defer fTemp2.Close()
	sTemp2 := bufio.NewScanner(fTemp2)
	sTemp2.Split(bufio.ScanLines)

	fTemp1, _ := os.Open("/Users/jrush/.quickBlocks/known_abis/uniq_sigs.tab")
	defer fTemp1.Close()
	sTemp1 := bufio.NewScanner(fTemp1)
	sTemp1.Split(bufio.ScanLines)

	temp2s := []string{""}
	for sTemp2.Scan() {
		temp2s = append(temp2s, sTemp2.Text())
	}

	for sTemp1.Scan() {
		for _, temp2 := range temp2s {
			wait.Add(1)
			_ = checkOne.Invoke(string(fmt.Sprintf("%s(%s)", temp2, sTemp1.Text())))
		}
	}
	// Wait till all threads finished at the end of program
	defer checkOne.Release()
	defer wait.Wait()
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "findSig",
	Short: "Do a brute force search for 4byte and event topics on cross product of names and signatures",
	Long: `
Description:
  This internal-use-only component searches the cross product of
  function and event names with likely signatures looking to match
  previously unmatched four bytes and event topics.`,
	Version: "GHC-TrueBlocks, LLC//0.8.1-alpha",
	Run: func(cmd *cobra.Command, args []string) {
		scrapeBlocks()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init Initalize options
func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.AddConfigPath(home + "/.quickBlocks")
	viper.SetConfigName("quickBlocks")
	viper.SetConfigType("toml")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.SetEnvPrefix("TB")
	viper.SetEnvKeyReplacer(strings.NewReplacer("SETTINGS.", ""))
	viper.AutomaticEnv()
}
