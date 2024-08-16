param (
    [Parameter(Mandatory=$true, Position=0, ValueFromPipeline=$true)]
    [string]$TargetPass
)

# Split the password into two parts
$firstPart = $TargetPass.Substring(0, [math]::Floor($TargetPass.Length / 2))
$secondPart = $TargetPass.Substring([math]::Floor($TargetPass.Length / 2))

# Base64 encode the parts
$encodedFirstPart = [Convert]::ToBase64String([Text.Encoding]::UTF8.GetBytes($firstPart))
$encodedSecondPart = [Convert]::ToBase64String([Text.Encoding]::UTF8.GetBytes($secondPart))

# Concatenate the parts with the specific pattern and base64 encode the whole string
$concatenated = "passM66QFT_$encodedFirstPart" + "_s5rS58O0O3ML_" + "$encodedSecondPart" + "_RendPassTS5S"
$encodedPassword = [Convert]::ToBase64String([Text.Encoding]::UTF8.GetBytes($concatenated))

# Print the final encoded password
Write-Output $encodedPassword