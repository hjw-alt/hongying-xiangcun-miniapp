$securePassword = Read-Host -Prompt "请输入 MySQL root 密码" -AsSecureString
$pointer = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($securePassword)
$password = [Runtime.InteropServices.Marshal]::PtrToStringBSTR($pointer)
[Runtime.InteropServices.Marshal]::ZeroFreeBSTR($pointer)

$encodedPassword = [Uri]::EscapeDataString($password)
$env:MYSQL_DSN = "root:$encodedPassword@tcp(127.0.0.1:3306)/shadowart_xiangtan?charset=utf8mb4&parseTime=true&loc=Local"

Write-Host "正在启动影绘湘潭后端..."
go run .
