$dir = Split-Path -Parent $MyInvocation.MyCommand.Path
$cert_name = "local-multi"

if (-not (Get-Command openssl -ErrorAction SilentlyContinue)) {
    Write-Host "Установи OpenSSL (choco install openssl)"
    exit 1
}

Push-Location $dir

# 🔐 Конфиг с несколькими доменами в SAN
$conf = @"
[req]
default_bits = 2048
prompt = no
default_md = sha256
distinguished_name = dn
x509_extensions = v3_req

[dn]
CN = localhost
O = LocalDev
C = RU

[v3_req]
subjectAltName = @alt_names
basicConstraints = CA:FALSE
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth

[alt_names]
DNS.1 = localhost
DNS.2 = shop.local
DNS.3 = local.hardwarestore.com
IP.1 = 127.0.0.1
IP.2 = ::1
"@

$conf | Out-File -FilePath "${cert_name}.cnf" -Encoding ASCII

# Генерация сертификата
openssl req -x509 -nodes -days 365 -newkey rsa:2048 `
    -keyout "${cert_name}.key" `
    -out "${cert_name}.crt" `
    -config "${cert_name}.cnf" `
    -extensions v3_req

Remove-Item "${cert_name}.cnf" -ErrorAction SilentlyContinue
Pop-Location

Write-Host "✅ Созданы: $dir\$cert_name.crt и $dir\$cert_name.key"
Write-Host "🔐 Установите $cert_name.crt в 'Доверенные корневые центры сертификации'"