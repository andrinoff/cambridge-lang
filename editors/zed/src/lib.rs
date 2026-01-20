use std::fs;
use zed_extension_api::{self as zed, Result};

struct CambridgeExtension {
    cached_binary_path: Option<String>,
}

impl CambridgeExtension {
    fn language_server_binary_path(
        &mut self,
        language_server_id: &zed::LanguageServerId,
        worktree: &zed::Worktree,
    ) -> Result<String> {
        // 1. Check if we already have the path cached in memory
        if let Some(path) = &self.cached_binary_path {
            if fs::metadata(path).map(|m| m.is_file()).unwrap_or(false) {
                return Ok(path.clone());
            }
        }

        // 2. Check if the LSP is already downloaded in the extension's support directory
        // Zed gives every extension a writable folder for this exact purpose.
        zed::set_language_server_installation_status(
            language_server_id,
            &zed::LanguageServerInstallationStatus::CheckingForUpdate,
        );

        let binary_name = "cambridge-lsp"; // Name of the file on disk
        let binary_path = format!("./{}", binary_name); // Path relative to the support dir

        if !fs::metadata(&binary_path)
            .map(|m| m.is_file())
            .unwrap_or(false)
        {
            // 3. DOWNLOAD IT if missing
            zed::set_language_server_installation_status(
                language_server_id,
                &zed::LanguageServerInstallationStatus::Downloading,
            );

            // Determine the user's OS and Architecture
            let (platform, arch) = zed::current_platform();

            // Construct the download URL based on the platform
            let download_url = match (platform, arch) {
                (zed::Os::Mac, zed::Architecture::Aarch64) =>
                    "https://github.com/andrinoff/cambridge-lang/releases/download/v0.0.1/cambridge-lsp-macos-arm64",
                (zed::Os::Mac, zed::Architecture::X8664) =>
                    "https://github.com/andrinoff/cambridge-lang/releases/download/v0.0.1/cambridge-lsp-macos-intel",
                (zed::Os::Linux, _) =>
                    "https://github.com/andrinoff/cambridge-lang/releases/download/v0.0.1/cambridge-lsp-linux",
                (zed::Os::Windows, _) =>
                    "https://github.com/andrinoff/cambridge-lang/releases/download/v0.0.1/cambridge-lsp.exe",
                _ => return Err("Unsupported platform".into()),
            };

            // Download the file
            zed::download_file(
                &download_url,
                &binary_path,
                zed::DownloadedFileType::Uncompressed, // Or Gzip/Zip if you compress it
            )
            .map_err(|e| format!("Failed to download LSP: {}", e))?;

            // Make it executable (Unix only)
            zed::make_file_executable(&binary_path)?;
        }

        self.cached_binary_path = Some(binary_path.clone());

        zed::set_language_server_installation_status(
            language_server_id,
            &zed::LanguageServerInstallationStatus::None,
        );

        Ok(binary_path)
    }
}

impl zed::Extension for CambridgeExtension {
    fn new() -> Self {
        Self {
            cached_binary_path: None,
        }
    }

    fn language_server_command(
        &mut self,
        language_server_id: &zed::LanguageServerId,
        worktree: &zed::Worktree,
    ) -> Result<zed::Command> {
        let path = self.language_server_binary_path(language_server_id, worktree)?;

        Ok(zed::Command {
            command: path,
            args: vec![],
            env: Default::default(),
        })
    }
}

zed::register_extension!(CambridgeExtension);
