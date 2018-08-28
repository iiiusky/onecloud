package shell

import (
	"fmt"

	"yunion.io/x/onecloud/pkg/util/aliyun"
	"yunion.io/x/onecloud/pkg/util/shellutils"
)

func init() {
	type InstanceListOptions struct {
		Id     []string `help:"IDs of instances to show"`
		Zone   string   `help:"Zone ID"`
		Limit  int      `help:"page size"`
		Offset int      `help:"page offset"`
	}
	shellutils.R(&InstanceListOptions{}, "instance-list", "List intances", func(cli *aliyun.SRegion, args *InstanceListOptions) error {
		instances, total, e := cli.GetInstances(args.Zone, args.Id, args.Offset, args.Limit)
		if e != nil {
			return e
		}
		printList(instances, total, args.Offset, args.Limit, []string{})
		return nil
	})

	type InstanceCrateOptions struct {
		NAME      string `help:"name of instance"`
		IMAGE     string `help:"image ID"`
		CPU       int    `help:"CPU count"`
		MEMORYGB  int    `help:"MemoryGB"`
		Disk      []int  `help:"Data disk sizes int GB"`
		STORAGE   string `help:"Storage type"`
		VSWITCH   string `help:"Vswitch ID"`
		PASSWD    string `help:"password"`
		PublicKey string `help:"PublicKey"`
	}
	shellutils.R(&InstanceCrateOptions{}, "instance-create", "Create a instance", func(cli *aliyun.SRegion, args *InstanceCrateOptions) error {
		instance, e := cli.CreateInstanceSimple(args.NAME, args.IMAGE, args.CPU, args.MEMORYGB, args.STORAGE, args.Disk, args.VSWITCH, args.PASSWD, args.PublicKey)
		if e != nil {
			return e
		}
		printObject(instance)
		return nil
	})

	type InstanceOperationOptions struct {
		ID string `help:"instance ID"`
	}
	shellutils.R(&InstanceOperationOptions{}, "instance-start", "Start a instance", func(cli *aliyun.SRegion, args *InstanceOperationOptions) error {
		err := cli.StartVM(args.ID)
		if err != nil {
			return err
		}
		return nil
	})

	shellutils.R(&InstanceOperationOptions{}, "instance-vnc", "Get a instance VNC url", func(cli *aliyun.SRegion, args *InstanceOperationOptions) error {
		url, err := cli.GetInstanceVNCUrl(args.ID)
		if err != nil {
			return err
		}
		fmt.Println(url)
		return nil
	})

	type InstanceStopOptions struct {
		ID    string `help:"instance ID"`
		Force bool   `help:"Force stop instance"`
	}
	shellutils.R(&InstanceStopOptions{}, "instance-stop", "Stop a instance", func(cli *aliyun.SRegion, args *InstanceStopOptions) error {
		err := cli.StopVM(args.ID, args.Force)
		if err != nil {
			return err
		}
		return nil
	})
	shellutils.R(&InstanceOperationOptions{}, "instance-delete", "Delete a instance", func(cli *aliyun.SRegion, args *InstanceOperationOptions) error {
		err := cli.DeleteVM(args.ID)
		if err != nil {
			return err
		}
		return nil
	})

	/*
		server-change-config 更改系统配置
		server-reset
	*/
	type InstanceDeployOptions struct {
		ID            string `help:"instance ID"`
		Name          string `help:"new instance name"`
		Hostname      string `help:"new hostname"`
		Keypair       string `help:"Keypair Name"`
		DeleteKeypair bool   `help:"Remove SSH keypair"`
		Password      string `help:"new password"`
		ResetPassword bool   `help:"Force reset password"`
		Description   string `help:"new instances description"`
	}

	shellutils.R(&InstanceDeployOptions{}, "instance-deploy", "Deploy keypair/password to a stopped virtual server", func(cli *aliyun.SRegion, args *InstanceDeployOptions) error {
		err := cli.DeployVM(args.ID, args.Name, args.Password, args.Keypair, args.ResetPassword, args.DeleteKeypair, args.Description)
		if err != nil {
			return err
		}
		return nil
	})

	type InstanceRebuildRootOptions struct {
		ID    string `help:"instance ID"`
		Image string `help:"Image ID"`
	}

	shellutils.R(&InstanceRebuildRootOptions{}, "instance-rebuild-root", "Reinstall virtual server system image", func(cli *aliyun.SRegion, args *InstanceRebuildRootOptions) error {
		err := cli.ReplaceSystemDisk(args.ID, args.Image)
		if err != nil {
			return err
		}
		return nil
	})

	type InstanceChangeConfigOptions struct {
		ID   string `help:"instance ID"`
		Ncpu int    `help:"number of CPU"`
		Vmem int    `help:"MiB of memory"`
		Disk []int  `help:"Data disk sizes int GB"`
	}

	shellutils.R(&InstanceChangeConfigOptions{}, "instance-change-config", "Deploy keypair/password to a stopped virtual server", func(cli *aliyun.SRegion, args *InstanceChangeConfigOptions) error {
		instance, e := cli.GetInstance(args.ID)
		if e != nil {
			return e
		}

		// todo : add create disks
		err := cli.ChangeVMConfig(instance.ZoneId, args.ID, args.Ncpu, args.Vmem, nil)
		if err != nil {
			return err
		}
		return nil
	})

	type InstanceUpdatePasswordOptions struct {
		ID     string `help:"Instance ID"`
		PASSWD string `help:"new password"`
	}
	shellutils.R(&InstanceUpdatePasswordOptions{}, "instance-update-password", "Update instance password", func(cli *aliyun.SRegion, args *InstanceUpdatePasswordOptions) error {
		err := cli.UpdateInstancePassword(args.ID, args.PASSWD)
		return err
	})
}
