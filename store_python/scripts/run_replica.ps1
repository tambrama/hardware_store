# Запуск дополнительного инстанса бэкенда (read-only, для балансировки)
# Пример: .\scripts\run_replica.ps1 -Port 5001 -DatabaseUri "postgresql://shop_readonly:readonly@localhost:5432/shopapi"

param(
    [int]$Port = 5001,
    [string]$DatabaseUri
)

$env:PORT = $Port
if ($DatabaseUri) { $env:DATABASE_URI = $DatabaseUri }
$root = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
Set-Location $root
python run.py
