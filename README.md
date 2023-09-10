# autoproofcli

autoproofcli is a Golang project that provides a command-line interface for automatic code & content protection.

## Installation

### Using pre-build binaries

#### Linux

You can download and install the latest version of your utility `autoproofcli` from GitHub Release:

1. Go to the `autoproofcli` repository page on GitHub: [github.com/autoproof/cli](https://github.com/autoproof/cli).
2. In the ["Releases"](https://github.com/Autoproof/cli/releases) section, find the latest released version. 
   Typically, it will be displayed at the top of the page.
3. Locate the release header that contains the version of your utility. For example, "Release v1.0.0."
4. Scroll down to the "Assets" section. You will see a list of available archives for download.
5. Find the archive that corresponds to your operating system. If you're using Linux, it might be a file 
   with the `.tar.gz` extension. Click on it to start the download.
6. After the archive has finished downloading, open a terminal on your system.
7. Extract the downloaded archive using the following command, replacing `<archive_filename.tar.gz>` with the 
   actual name of the downloaded file:
   ```bash
   tar -xzvf <archive_filename.tar.gz> 
   ```
8. Now you have access to the `autoproofcli` executable. You can either run it from the current directory 
   or copy it to a directory in your PATH variable to run `autoproofcli` from any directory on your system. 
9. To run `autoproofcli` from the current directory, use the following command:
    ```bash
    ./autoproofcli --help
    ```
10. To copy `autoproofcli` to a directory in your PATH (e.g., `/usr/local/bin/`), execute the following command 
    with administrator privileges:
    ```bash
    sudo cp autoproofcli /usr/local/bin/
    ```
11. You can now use `autoproofcli` from any directory on your system by simply typing its name in the command line:
    ```bash
    autoproofcli
    ```

#### Windows

You can download and install the latest version of your utility `autoproofcli` from GitHub Releases on Windows:

1. Go to the `autoproofcli` repository page on GitHub: [github.com/autoproof/cli](https://github.com/autoproof/cli).
2. In the ["Releases"](https://github.com/Autoproof/cli/releases) section, find the latest released version.
   Typically, it will be displayed at the top of the page.
3. Locate the release header that contains the version of your utility. For example, "Release v1.0.0."
4. Scroll down to the "Assets" section. You will see a list of available archives for download.
5. Find the archive that corresponds to your operating system. Since you are using Windows, it will likely 
   be a zip archive containing executable files (`.zip`). Click on it to start the download.
6. After the archive has finished downloading, locate the downloaded zip file in your downloads folder or
   the directory where you saved it.
7. Extract the contents of the zip archive to a location of your choice. You can do this by right-clicking 
   the zip file and selecting "Extract All."
8. Navigate to the directory where you extracted the contents of the zip archive using the File Explorer.
9. You will find the `autoproofcli.exe` executable file in the extracted folder. This is the main executable of 
   your utility.
10. You can now run `autoproofcli.exe` by double-clicking on it. If you want to use it from the command prompt, 
    open the Command Prompt or PowerShell and navigate to the folder where `autoproofcli.exe` is located.
11. To run `autoproofcli.exe` from the command prompt or PowerShell, use the following command:

    ```
    autoproofcli.exe
    ```
12. You have now installed and can use the latest version of your `autoproofcli` utility on your Windows system.

### Using the source code

To install autoproofcli, make sure you have Golang installed on your system. 
If you don't have Golang installed, you can download it from the official Golang website: https://golang.org/dl/

Once you have Golang installed, follow the steps below to build the autoproofcli project:

1. Clone the autoproof CLI repository:

   ```
   git clone https://github.com/autoproof/cli.git
   ```

2. Change to the project directory:

   ```
   cd cli
   ```

3. Build the project using the `go build` command:

   ```
   go build -o autoproofcli
   ```

   This command will compile the project and create an executable binary called `autoproofcli`.


## Usage

After you have built or installed the `autoproofcli` binary, you can use it by running the executable from the 
command line. Here's an example of how to use it:

```
autoproofcli command [arguments]
```

Replace `command` with the specific command you want to execute and provide any required arguments. You can find 
information about available commands and their usage in the project documentation or by running `autoproofcli --help`.


## CI / CD Integration

### GitHub Actions

This GitHub Action Workflow is meant to automatically generate a snapshot using the Autoproof tool whenever 
you push to the main branch of your repository:

```yaml
name: Autoproof Snapshot
on:
  push:
    branches:
      - main
jobs:
  autoproof:
    runs-on: ubuntu-latest
    # NOTE: Instead of the "latest" tag, we recommend using the stable version tag of the autoproof/cli image 
    # to have explicit control over updates (to avoid breaking backwards compatibility).
    container: ghcr.io/autoproof/cli:latest
    steps:
      - name: Checkout source code
        uses: actions/checkout@v3
        with:
          fetch-depth: 0
      - name: Autoproof snapshot
        env:
          AUTOPROOF_APIKEY: ${{ secrets.AUTOPROOF_APIKEY }}
        run: |
          autoproofcli snapshot -m "GHA on ${{ github.repository }}@${{ github.sha }}: ${{ github.event.head_commit.message }}"
```

Security Note:

Make sure your Autoproof API Key (AUTOPROOF_APIKEY) is stored in your repository's secrets for confidentiality.
