# 输出当前可分享给同一局域网内其他人的访问地址
$ips = Get-NetIPAddress -AddressFamily IPv4 |
  Where-Object { $_.IPAddress -notlike '127.*' -and $_.IPAddress -notlike '169.254.*' } |
  Sort-Object -Property InterfaceAlias

$running = Get-NetTCPConnection -LocalPort 8080 -State Listen -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "红映湘村 · 共享地址" -ForegroundColor Red
Write-Host "-----------------------------------------"

if (-not $running) {
  Write-Host "后端服务未启动：请先在 backend 目录执行 go run ." -ForegroundColor Yellow
}

if (-not $ips) {
  Write-Host "未检测到局域网地址，请确认已连接 WiFi 或热点。" -ForegroundColor Yellow
}

foreach ($ip in $ips) {
  Write-Host "网卡：$($ip.InterfaceAlias)"
  Write-Host "  管理后台：http://$($ip.IPAddress):8080/admin/" -ForegroundColor Green
  Write-Host "  接口地址：http://$($ip.IPAddress):8080/api"
  Write-Host "  账号密码：hjw / 123456"
  Write-Host ""
}

Write-Host "注意：对方需要连接同一个 WiFi 或热点；换网络后地址会变化。" -ForegroundColor DarkGray
