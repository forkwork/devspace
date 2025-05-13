use tauri::AppHandle;
use crate::resource_watcher::ProInstance;

use super::{
    config::{CommandConfig, DevspaceCommandConfig, DevspaceCommandError},
    constants::{DEVSPACE_BINARY_NAME, DEVSPACE_COMMAND_LIST, DEVSPACE_COMMAND_PRO, FLAG_OUTPUT_JSON},
};

pub struct ListProInstancesCommand {}
impl ListProInstancesCommand {
    pub fn new() -> Self {
        ListProInstancesCommand {}
    }

    fn deserialize(&self, d: Vec<u8>) -> Result<Vec<ProInstance>, DevspaceCommandError> {
        serde_json::from_slice(&d).map_err(DevspaceCommandError::Parse)
    }
}
impl DevspaceCommandConfig<Vec<ProInstance>> for ListProInstancesCommand {
    fn config(&self) -> CommandConfig {
        CommandConfig {
            binary_name: DEVSPACE_BINARY_NAME,
            args: vec![DEVSPACE_COMMAND_PRO, DEVSPACE_COMMAND_LIST, FLAG_OUTPUT_JSON],
        }
    }

    fn exec_blocking(self, app_handle: &AppHandle) -> Result<Vec<ProInstance>, DevspaceCommandError> {
        let cmd = self.new_command(app_handle)?;

        let output = tauri::async_runtime::block_on(async move { cmd.output().await })
            .map_err(|_| DevspaceCommandError::Output)?;

        self.deserialize(output.stdout)
    }
}
impl ListProInstancesCommand {
    pub async fn exec(
        self,
        app_handle: &AppHandle,
    ) -> Result<Vec<ProInstance>, DevspaceCommandError> {
        let cmd = self.new_command(app_handle)?;

        let output = cmd.output().await.map_err(|_| DevspaceCommandError::Output)?;

        self.deserialize(output.stdout)
    }
}
