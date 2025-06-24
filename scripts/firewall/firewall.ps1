$port = 60011
$ruleName = "fileshare"
$description = "允许文件共享服务通过 TCP 端口 $port 进行入站连接。"

Write-Host "检查是否存在名为 '$ruleName' 的防火墙规则..."

$existingRule = Get-NetFirewallRule -DisplayName $ruleName -ErrorAction SilentlyContinue

if ($existingRule) {
    Write-Host "防火墙规则 '$ruleName' 已存在。删除旧规则以便更新..."
    Remove-NetFirewallRule -DisplayName $ruleName -Confirm:$false
    Write-Host "旧规则已删除。"
} else {
    Write-Host "防火墙规则 '$ruleName' 不存在。将创建新规则。"
}

Write-Host "正在创建新的入站防火墙规则..."
try {
    New-NetFirewallRule -DisplayName $ruleName `
                        -Description $description `
                        -Direction Inbound `
                        -Action Allow `
                        -Protocol TCP `
                        -LocalPort $port `
                        -Profile Any `
                        -Enabled True

    Write-Host "成功创建防火墙规则: '$ruleName' (端口: $port/TCP)"
}
catch {
    Write-Host "创建防火墙规则时发生错误: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "请确保您以管理员身份运行此脚本。" -ForegroundColor Yellow
}

Write-Host "操作完成。"