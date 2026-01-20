use std::fs;
use zed_extension_api::{self as zed, Result};

struct CambridgeExtension {
    cached_binary_path: Option<String>,
}

impl CambridgeExtension {
    fn language_server_binary_path(
        &mut self,
        language_server_id: &zed::LanguageServerId,
        _worktree: &zed::Worktree,
    ) -> Result<String> {
        // 1. Check if we already have the path cached in memory
        if let Some(path) = &self.cached_binary_path {
            if fs::metadata(path).map(|m| m.is_file()).unwrap_or(false) {
                return Ok(path.clone());
            }
        }

        zed::set_language_server_installation_status(
            language_server_id,
            &zed::LanguageServerInstallationStatus::CheckingForUpdate,
        );

        let binary_name = "cambridge-lsp";
        // This path is relative to the extension's installation folder (support dir)
        let binary_path = format!("./{}", binary_name);

        if !fs::metadata(&binary_path)
            .map(|m| m.is_file())
            .unwrap_or(false)
        {
            zed::set_language_server_installation_status(
                language_server_id,
                &zed::LanguageServerInstallationStatus::Downloading,
            );

            let (platform, arch) = zed::current_platform();

            // When you are ready to release v0.0.2 with the new features,
            // update these URLs to point to the new assets.
            let download_url = match (platform, arch) {
                (zed::Os::Mac, zed::Architecture::Aarch64) =>
                    "https://github.com/andrinoff/cambridge-lang/releases/download/v0.1.0/cambridge-lsp-macos-arm64",
                (zed::Os::Mac, zed::Architecture::X8664) =>
                    "https://github.com/andrinoff/cambridge-lang/releases/download/v0.1.0/cambridge-lsp-macos-intel",
                (zed::Os::Linux, _) =>
                    "https://github.com/andrinoff/cambridge-lang/releases/download/v0.1.0/cambridge-lsp-linux",
                (zed::Os::Windows, _) =>
                    "https://github.com/andrinoff/cambridge-lang/releases/download/v0.1.0/cambridge-lsp.exe",
                _ => return Err("Unsupported platform".into()),
            };

            zed::download_file(
                &download_url,
                &binary_path,
                zed::DownloadedFileType::Uncompressed,
            )
            .map_err(|e| format!("Failed to download LSP: {}", e))?;

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
