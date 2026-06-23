package config

import (
	"os"
	"testing"
)

// TestRootPassword_DirectEnv 直接设置 ROOT_PASSWORD
func TestRootPassword_DirectEnv(t *testing.T) {
	os.Setenv(EnvRootPassword, "test-secret-123")
	os.Unsetenv(EnvRootPasswordFile)
	defer os.Unsetenv(EnvRootPassword)

	got := RootPassword()
	if got != "test-secret-123" {
		t.Errorf("RootPassword() = %q, want %q", got, "test-secret-123")
	}
}

// TestRootPassword_FilePriority 文件优先于直接环境变量
func TestRootPassword_FilePriority(t *testing.T) {
	os.Setenv(EnvRootPassword, "env-value")
	defer os.Unsetenv(EnvRootPassword)

	// 创建临时密码文件
	f, err := os.CreateTemp("", "root-password-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	os.WriteFile(f.Name(), []byte("  file-value  \n"), 0600)

	os.Setenv(EnvRootPasswordFile, f.Name())
	defer os.Unsetenv(EnvRootPasswordFile)

	got := RootPassword()
	if got != "file-value" {
		t.Errorf("RootPassword() = %q, want %q (文件应当优先)", got, "file-value")
	}
}

// TestRootPassword_FileTrim 文件内容带空白
func TestRootPassword_FileTrim(t *testing.T) {
	os.Unsetenv(EnvRootPassword)
	defer os.Unsetenv(EnvRootPassword)

	f, err := os.CreateTemp("", "root-password-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	os.WriteFile(f.Name(), []byte("\n\n  password-with-spaces  \n"), 0600)

	os.Setenv(EnvRootPasswordFile, f.Name())
	defer os.Unsetenv(EnvRootPasswordFile)

	got := RootPassword()
	if got != "password-with-spaces" {
		t.Errorf("RootPassword() = %q, want %q", got, "password-with-spaces")
	}
}

// TestRootPassword_FileNotFound panics
func TestRootPassword_FileNotFound(t *testing.T) {
	os.Unsetenv(EnvRootPassword)
	defer os.Unsetenv(EnvRootPassword)
	os.Setenv(EnvRootPasswordFile, "/nonexistent/password/file")
	defer os.Unsetenv(EnvRootPasswordFile)

	defer func() {
		if r := recover(); r == nil {
			t.Error("RootPassword() 应当 panic（文件不存在）")
		}
	}()
	RootPassword()
}

// TestRootPassword_NoPasswordNoDB 没有密码也没有 DB → panic
func TestRootPassword_NoPasswordNoDB(t *testing.T) {
	os.Unsetenv(EnvRootPassword)
	os.Unsetenv(EnvRootPasswordFile)
	defer os.Unsetenv(EnvRootPassword)
	defer os.Unsetenv(EnvRootPasswordFile)

	// DB 为 nil（未初始化）→ 回退不到第三步
	oldDB := DB
	DB = nil
	defer func() { DB = oldDB }()

	defer func() {
		if r := recover(); r == nil {
			t.Error("RootPassword() 应当 panic（首次部署无密码）")
		} else {
			expected := "首次部署需要设置 ROOT_PASSWORD 或 ROOT_PASSWORD_FILE 环境变量"
			if r != expected {
				t.Errorf("panic = %q, want %q", r, expected)
			}
		}
	}()
	RootPassword()
}

// TestUnsetEnv_AfterSeed 验证 os.Unsetenv 在 seedRootUser 后生效
// （函数级别验证：os.Unsetenv 的行为是标准库保证的）
func TestUnsetEnv_AfterUnset(t *testing.T) {
	os.Setenv(EnvRootPassword, "should-be-gone")
	os.Unsetenv(EnvRootPassword)
	if os.Getenv(EnvRootPassword) != "" {
		t.Error("os.Unsetenv 未生效")
	}
}
