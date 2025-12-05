package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/sys/unix"
)

type Config struct {
	Name   string       `json:"name"`
	Cgroup CgroupConfig `json:"cgroup"`
	Rootfs RootfsConfig `json:"rootfs"`
}

func main() {
	// このgoroutineが実行されるOSスレッドを1つに定め、固定
	//  Namespaceやcgroupの設定を正しく行うため
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// エントリーポイント以外の設定の読み込み
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	var c Config
	if err := json.Unmarshal(configFile, &c); err != nil {
		log.Fatalln(err)
	}

	// 指定されたコマンドの実行
	switch os.Args[1] {
	case "run":
		if err := run(c, os.Args[2:]); err != nil {
			log.Fatalln(err)
		}

	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}
}

// runコマンド
func run(c Config, entryPoint []string) error {
	// Namespaceの設定
	if err := SetupNamespace(); err != nil {
		return err
	}

	// cgroupの設定
	if err := SetupCgroup(c.Name, os.Getpid(), c.Cgroup); err != nil {
		return err
	}

	// rootfsの設定
	_ = unix.Unshare(unix.CLONE_NEWNS) // rootfsで使うので、Namespace系の処理だが仮置き
	if err := SetupRootfs(c.Rootfs); err != nil {
		return err
	}

	// コンテナ(仮)内でエントリーポイントを実行
	cmd := exec.Command(entryPoint[0], entryPoint[1:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}
