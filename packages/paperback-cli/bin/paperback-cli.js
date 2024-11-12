#!/usr/bin/env node

function getBinaryPath() {
  // Lookup table for all platforms and binary distribution packages
  const BINARY_DISTRIBUTION_PACKAGES = {
    'darwin-arm': 'paperback-cli-darwin-arm64',
    'darwin-x64': 'paperback-cli-darwin-amd64',
    'linux-x64': 'paperback-cli-linux-amd64',
    'win32-x64': 'paperback-cli-windows-amd64',
  }

  // Windows binaries end with .exe so we need to special case them.
  const binaryName = process.platform === 'win32' ? 'my-binary.exe' : 'my-binary'

  // Determine package name for this platform
  const platformSpecificPackageName =
    BINARY_DISTRIBUTION_PACKAGES[`${process.platform}-${process.arch}`]

  try {
    // Resolving will fail if the optionalDependency was not installed
    return require.resolve(`${platformSpecificPackageName}/bin/${binaryName}`)
  } catch (e) {
    return require('path').join(__dirname, '..', binaryName)
  }
}

require('child_process').execFileSync(getBinaryPath(), process.argv.slice(2), {
  stdio: 'inherit',
})