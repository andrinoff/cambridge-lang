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
        // 1. If we already found the binary, return it
        if let Some(path) = &self.cached_binary_path {
            if std::fs::metadata(path)
                .map(|m| m.is_file())
                .unwrap_or(false)
            {
                return Ok(path.clone());
            }
        }

        // 2. Look for 'cambridge-lsp' in the user's PATH
        if let Some(path) = worktree.which("cambridge-lsp") {
            self.cached_binary_path = Some(path.clone());
            return Ok(path);
        }

        // 3. Fallback: If not in PATH, you might look in the project folder
        // For development, we assume the user has built it and put it in PATH
        Err(format!(
            "Protocol error: 'cambridge-lsp' binary not found. \
            Please ensure you have built the language server and added it to your PATH."
        ))
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
