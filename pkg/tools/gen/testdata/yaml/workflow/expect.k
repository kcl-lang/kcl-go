"""
This file was generated by the KCL auto-gen tool. DO NOT EDIT.
Editing this file might prove futile when you re-run the KCL auto-gen generate command.
"""

{
    name = "Daily Release CI"
    on = {
        schedule = [
            {
                cron = "0 0 * * *"
            }
        ]
    }
    permissions = {
        contents = "write"
    }
    jobs = {
        "macos-release" = {
            name = "Build and release on macos"
            "runs-on" = "macos-11"
            steps = [
                {
                    name = "Git checkout"
                    uses = "actions/checkout@v2"
                    with = {
                        submodules = "true"
                    }
                }
                {
                    run = "clang --version"
                }
                {
                    run = "cargo --version"
                }
                {
                    run = "rustc --print sysroot"
                }
                {
                    name = "Delete rust cargo"
                    run = "rm -rf /root/.cargo/bin"
                    shell = "bash"
                }
                {
                    name = "Install LLVM 12"
                    run = "brew install llvm@12"
                    shell = "bash"
                }
                {
                    name = "Install rust nightly toolchain"
                    uses = "actions-rs/toolchain@v1"
                    with = {
                        toolchain = 1.73
                        override = True
                        components = "clippy, rustfmt"
                    }
                }
                {
                    name = "Build KCL"
                    run = r"""export PATH=$PATH:$PWD/../_build/dist/Darwin/kclvm/bin:/usr/local/opt/llvm@12/bin && make build"""
                    shell = "bash"
                }
                {
                    uses = "actions/upload-artifact@v3"
                    with = {
                        name = "kcl-darwin-amd64-nightly"
                        path = "_build/dist/Darwin/kclvm"
                    }
                }
            ]
        }
        "linux-release" = {
            name = "Build and release on linux"
            "runs-on" = "ubuntu-latest"
            steps = [
                {
                    name = "Check out code"
                    uses = "actions/checkout@v3"
                    with = {
                        submodules = "true"
                    }
                }
                {
                    name = "Install LLVM"
                    shell = "bash"
                    run = r"""sudo apt-get update
sudo apt-get install -y git wget curl make
sudo apt-get install -y clang-12 lld-12
sudo ln -sf /usr/bin/clang-12 /usr/bin/clang
"""
                }
                {
                    name = "Install Rust Toolchain"
                    uses = "actions-rs/toolchain@v1"
                    with = {
                        toolchain = 1.73
                        override = True
                        components = "clippy, rustfmt"
                    }
                }
                {
                    name = "Build KCL"
                    run = "make build"
                    shell = "bash"
                }
                {
                    uses = "actions/upload-artifact@v3"
                    with = {
                        name = "kcl-linux-amd64-nightly"
                        path = "_build/dist/ubuntu/kclvm"
                    }
                }
                {
                    name = "Copy Dockerfile to the current work directory"
                    run = "cp scripts/docker/kcl/Dockerfile ."
                    shell = "bash"
                }
                {
                    name = "Login to Docker Hub"
                    uses = "docker/login-action@v2"
                    with = {
                        username = r"""${{ secrets.DOCKER_USERNAME }}"""
                        password = r"""${{ secrets.DOCKER_PASSWORD }}"""
                    }
                }
                {
                    name = "Extract metadata (tags, labels) for Docker"
                    id = "meta"
                    uses = "docker/metadata-action@v4"
                    with = {
                        images = "kcllang/kclvm"
                    }
                }
                {
                    name = "Build and push Docker image"
                    uses = "docker/build-push-action@v3"
                    with = {
                        context = "."
                        push = r"""${{ github.event_name != 'pull_request' }}"""
                        tags = r"""${{ steps.meta.outputs.tags }}"""
                        labels = r"""${{ steps.meta.outputs.labels }}"""
                    }
                }
            ]
        }
        "windows-release" = {
            name = "Build and release on windows"
            "runs-on" = "windows-latest"
            env = {
                LLVM_SYS_120_PREFIX = "C:/LLVM"
                KCLVM_CLANG = "C:/LLVM/bin/clang.exe"
            }
            steps = [
                {
                    name = "Git checkout"
                    uses = "actions/checkout@v2"
                    with = {
                        submodules = "true"
                    }
                }
                {
                    uses = "ilammy/msvc-dev-cmd@v1"
                }
                {
                    run = "clang --version"
                }
                {
                    run = "cargo --version"
                }
                {
                    run = "Invoke-WebRequest -Uri https://github.com/kcl-lang/llvm-package-windows/releases/download/v12.0.1/LLVM-12.0.1-win64.7z -OutFile C:/LLVM-12.0.1-win64.7z"
                }
                {
                    run = "Get-FileHash -Algorithm MD5 C:/LLVM-12.0.1-win64.7z"
                }
                {
                    run = r"""7z x -y C:/LLVM-12.0.1-win64.7z -o"C:/LLVM" """
                }
                {
                    run = "Remove-Item C:/LLVM-12.0.1-win64.7z"
                }
                {
                    run = r"""echo "C:/LLVM/bin" | Out-File -FilePath $env:GITHUB_PATH -Encoding utf8 -Append"""
                }
                {
                    run = r""".\scripts\build-windows\build.ps1"""
                }
                {
                    run = r"""echo ";$(pwd)\scripts\build-windows\_output\kclvm-windows\bin" | Out-File -FilePath $env:GITHUB_PATH -Encoding utf8 -Append"""
                    "working-directory" = "."
                }
                {
                    uses = "actions/upload-artifact@v3"
                    with = {
                        name = "kcl-windows-nightly"
                        path = "scripts/build-windows/_output/kclvm-windows"
                    }
                }
            ]
        }
    }
}
