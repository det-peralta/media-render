[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://buymeacoffee.com/detperalta)

# VidFusion

VidFusion is a video processing tool that supports hardware-accelerated video conversion and concatenation.

## Usage
You can download the pre-built executable directly from the [GitHub Releases](https://github.com/your-repo/vidfusion/releases) page. Once downloaded:

1. Drop video files onto the executable or pass them as arguments.
2. Follow the prompts to configure options like audio removal and video conversion.
3. The output file will be saved as `output.mp4`.

## Build or Development
If you prefer to build VidFusion locally or develop it further, follow these steps:

### Requirements
- [FFmpeg](https://ffmpeg.org/) must be installed and available in your system's PATH.
- Go 1.20 or later.

### Building Locally
#### For Windows
```bash
GOOS=windows GOARCH=amd64 go build -o vidfusion_windows.exe vidfusion.go
```

#### For Linux
```bash
GOOS=linux GOARCH=amd64 go build -o vidfusion_linux vidfusion.go
```

## Setting up in GitHub Codespaces
You can use GitHub Codespaces to develop and test VidFusion in a cloud-based environment. Follow these steps:

1. Open the repository in a Codespace.
2. Ensure FFmpeg is installed in the Codespace:
   ```bash
   sudo apt update && sudo apt install -y ffmpeg
   ```
3. Build the project:
   ```bash
   go build -o vidfusion vidfusion.go
   ```
4. Run the application:
   ```bash
   ./vidfusion <video-files>
   ```

## Installing FFmpeg on Windows
To install FFmpeg on Windows using `winget`, run the following command:
```bash
winget install -e --id Gyan.FFmpeg
```
Ensure that FFmpeg is added to your system's PATH after installation.
