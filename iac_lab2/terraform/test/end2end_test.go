package test

import (
	"fmt"
	"io/ioutil"
	"flag"
	"testing"
	"database/sql"
	"crypto/tls"
	"crypto/x509"
	"os"
	// "reflect"


	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"golang.org/x/crypto/ssh"
	"github.com/go-sql-driver/mysql"
)

var folder = flag.String("folder", "", "Folder ID in Yandex.Cloud")
var sshKeyPath = flag.String("ssh-key-pass", "", "Private ssh key for access to virtual machines")
var dbSshKeyPath = flag.String("db-ssh-key-pass", "", "Private ssh key for access to db")

func TestEndToEndDeploymentScenario(t *testing.T) {
    fixtureFolder := "../"

    test_structure.RunTestStage(t, "setup", func() {
		terraformOptions := &terraform.Options{
			TerraformDir: fixtureFolder,

			Vars: map[string]interface{}{
			"yc_folder":    *folder,
		    },
			
	    }

		test_structure.SaveTerraformOptions(t, fixtureFolder, terraformOptions)

		terraform.InitAndApply(t, terraformOptions)
	})

	test_structure.RunTestStage(t, "validate", func() {
	    fmt.Println("Run some tests...")

	    terraformOptions := test_structure.LoadTerraformOptions(t, fixtureFolder)

        // test load balancer ip existing
	    loadbalancerIPAddress := terraform.Output(t, terraformOptions, "load_balancer_public_ip")

	    if loadbalancerIPAddress == "" {
			t.Fatal("Cannot retrieve the public IP address value for the load balancer.")
		}

		// test ssh connect
		vmLinuxPublicIPAddress := terraform.Output(t, terraformOptions, "vm_linux_public_ip_address")

		key, err := ioutil.ReadFile(*sshKeyPath)
		if err != nil {
			t.Fatalf("Unable to read private key: %v", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			t.Fatalf("Unable to parse private key: %v", err)
		}

		sshConfig := &ssh.ClientConfig{
			User: "ubuntu",
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		sshConnection, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", vmLinuxPublicIPAddress), sshConfig)
		if err != nil {
			t.Fatalf("Cannot establish SSH connection to vm-linux public IP address: %v", err)
		}

		defer sshConnection.Close()
        
		sshSession, err := sshConnection.NewSession()
		if err != nil {
			t.Fatalf("Cannot create SSH session to vm-linux public IP address: %v", err)
		}

		defer sshSession.Close()
        
		err = sshSession.Run(fmt.Sprintf("ping -c 1 8.8.8.8"))
		if err != nil {
			t.Fatalf("Cannot ping 8.8.8.8: %v", err)
		}

		// test db
		databaseHostFQDNs := terraform.OutputList(t, terraformOptions, "database_host_fqdn")
		dbuser := os.Getenv("DBUSER")
		dbpass := os.Getenv("DBPASS")
		net := "tcp"
		dbname := "db"
		port := 3306

		// fmt.Println(reflect.TypeOf(databaseHostFQDNs))
		// for _, host := range databaseHostFQDNs {
		// 	fmt.Printf(fmt.Sprintf("%s:%s@%s(%v:%d)/%s?tls=custom\n",	dbuser, dbpass, net, host, port, dbname))
		// }
		
		if dbuser == "" || dbpass == "" {
			t.Fatal("env DBUSER or DBPASS not found!!")
		}

		rootCertPool := x509.NewCertPool()
		dbkey, err := ioutil.ReadFile(*dbSshKeyPath)
		if err != nil {
			t.Fatalf("Unable to read private key dbSshKeyPath: %v", err)
		}

		if !rootCertPool.AppendCertsFromPEM(dbkey) {
			t.Fatalf("AppendCertsFromPEM err")
		}
		
		mysql.RegisterTLSConfig("custom", &tls.Config{
			RootCAs: rootCertPool,
		})

		for _, host := range databaseHostFQDNs {
			mysqlInfo := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?tls=custom",
				dbuser, dbpass, net, host, port, dbname)

			conn, err := sql.Open("mysql", mysqlInfo)
			if err != nil {
				t.Fatalf("ERROR! %s: %v", host, err)
			}
			
			fmt.Println(fmt.Sprintf("Connected! %s", host))
			defer conn.Close()

			q, err := conn.Query("SELECT version()")
			if err != nil {
				t.Fatalf("ERROR! %s: %v", host, err)
			}

			var res string

			for q.Next() {
				q.Scan(&res)
				fmt.Printf(fmt.Sprintf("MySql version %s: %s\n", host, res))
			}
		}
    })

	test_structure.RunTestStage(t, "teardown", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, fixtureFolder)
		terraform.Destroy(t, terraformOptions)
	})
}