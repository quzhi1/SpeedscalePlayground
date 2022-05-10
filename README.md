# SpeedscalePlayground
My online lab for speedscale

## Prerquisite
1. Have minikube up and running.
2. Install tilt: `curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash`
3. Go to https://app.speedscale.com/profile, and put your API key in .zshrc:
```bash
export SPEEDSCALE_API_KEY=<your-api-key>
```

## Running services
```bash
tilt up
```
