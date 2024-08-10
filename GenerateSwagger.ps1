# This powershell script will generate swagger api docs for this project

# make sure the swag command exists
if (!(Get-Command swag -ErrorAction SilentlyContinue)) {
    Write-Host "Installing swag command line"
    & go install github.com/swaggo/swag/cmd/swag@latest
}

# generate swagger api docs
Write-Host "Generating swagger api docs"
& swag init --output src/apiHandlers/swaggerHandlers
