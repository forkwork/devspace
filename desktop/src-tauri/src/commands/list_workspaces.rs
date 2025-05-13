use tauri::AppHandle;

use crate::resource_watcher::Workspace;

use super::{
    config::{CommandConfig, DevspaceCommandConfig, DevspaceCommandError},
    constants::{DEVSPACE_BINARY_NAME, DEVSPACE_COMMAND_LIST, FLAG_OUTPUT_JSON},
};

pub struct ListWorkspacesCommand {}
impl ListWorkspacesCommand {
    pub fn new() -> Self {
        ListWorkspacesCommand {}
    }

    fn deserialize(&self, d: Vec<u8>) -> Result<Vec<Workspace>, DevspaceCommandError> {
        serde_json::from_slice(&d).map_err(DevspaceCommandError::Parse)
    }
}
impl DevspaceCommandConfig<Vec<Workspace>> for ListWorkspacesCommand {
    fn config(&self) -> CommandConfig {
        CommandConfig {
            binary_name: DEVSPACE_BINARY_NAME,
            args: vec![DEVSPACE_COMMAND_LIST, FLAG_OUTPUT_JSON],
        }
    }

    fn exec_blocking(self, app_handle: &AppHandle) -> Result<Vec<Workspace>, DevspaceCommandError> {
        let cmd = self.new_command(app_handle)?;

        let output = tauri::async_runtime::block_on(async move { cmd.output().await })
            .map_err(|_| DevspaceCommandError::Output)?;

        self.deserialize(output.stdout)
    }
}

impl ListWorkspacesCommand {
    pub async fn exec(self, app_handle: &AppHandle) -> Result<Vec<Workspace>, DevspaceCommandError> {
        let cmd = self.new_command(app_handle)?;

        let output = cmd.output().await.map_err(|_| DevspaceCommandError::Output)?;

        self.deserialize(output.stdout)
    }
}
