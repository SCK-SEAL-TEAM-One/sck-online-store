"""
Helper library for downloading files from Selenium Grid remote browsers.
Uses Selenium 4's download API for remote file downloads.
"""

from robot.libraries.BuiltIn import BuiltIn
from selenium import webdriver
import os
import time


class DownloadHelper:
    """Library for handling file downloads from Selenium Grid."""

    ROBOT_LIBRARY_SCOPE = 'TEST'

    def __init__(self):
        self.selenium_lib = None

    def _get_selenium_library(self):
        """Get SeleniumLibrary instance."""
        if self.selenium_lib is None:
            self.selenium_lib = BuiltIn().get_library_instance('SeleniumLibrary')
        return self.selenium_lib

    def _get_driver(self):
        """Get the current WebDriver instance."""
        selenium_lib = self._get_selenium_library()
        return selenium_lib.driver

    def enable_download_in_headless_chrome(self, download_path="/home/seluser/downloads"):
        """
        Enable download in headless Chrome using CDP.
        Must be called after opening browser and before clicking download.
        """
        driver = self._get_driver()
        
        try:
            # Enable download using Chrome DevTools Protocol
            driver.execute_cdp_cmd("Page.setDownloadBehavior", {
                "behavior": "allow",
                "downloadPath": download_path
            })
            BuiltIn().log(f"Download enabled at: {download_path}", level='INFO')
        except Exception as e:
            BuiltIn().log(f"Could not enable downloads via CDP: {e}", level='WARN')

    def get_downloaded_files_from_remote(self, timeout=30):
        """
        Get list of downloaded files from remote browser.
        - For Selenium Grid: copy files from container to local directory using docker cp
        - For Local browser: check local directory directly
        
        Returns list of dictionaries with 'name' and 'path' keys.
        """
        import subprocess
        
        download_dir = BuiltIn().get_variable_value('${DOWNLOAD_DIR}')
        remote_hub_url = BuiltIn().get_variable_value('${REMOTE_HUB_URL}', default='')
        
        # Create download directory if it doesn't exist
        os.makedirs(download_dir, exist_ok=True)
        
        start_time = time.time()
        downloaded_files = []
        copied_files = set()
        
        BuiltIn().log(f"Waiting for downloads to complete...", level='INFO')
        BuiltIn().log(f"Download directory: {download_dir}", level='INFO')
        BuiltIn().log(f"Remote Hub URL: {remote_hub_url if remote_hub_url else 'Not set (Local mode)'}", level='INFO')
        
        # Check if running in remote (Selenium Grid) or local mode
        is_remote = remote_hub_url and remote_hub_url.strip() != ''
        
        while time.time() - start_time < timeout:
            try:
                elapsed = time.time() - start_time
                
                if is_remote:
                    # Grid mode: List files from container and copy them
                    result = subprocess.run(
                        ['docker', 'exec', 'node-chrome', 'find', '/home/seluser/downloads', '-name', '*.pdf', '-type', 'f'],
                        capture_output=True,
                        text=True,
                        timeout=5
                    )
                    
                    if result.returncode == 0:
                        container_files = [f.strip() for f in result.stdout.split('\n') if f.strip().endswith('.pdf')]
                        BuiltIn().log(f"[{elapsed:.1f}s] Checking container... found {len(container_files)} PDF(s)", level='DEBUG')
                        
                        if container_files:
                            # Copy files from container to local directory
                            for file_path in container_files:
                                file_name = os.path.basename(file_path)
                                
                                # Skip if already copied
                                if file_name in copied_files:
                                    continue
                                
                                # Copy file from container
                                local_path = os.path.join(download_dir, file_name)
                                copy_result = subprocess.run(
                                    ['docker', 'cp', f'node-chrome:{file_path}', local_path],
                                    capture_output=True,
                                    text=True,
                                    timeout=10
                                )
                                
                                if copy_result.returncode == 0 and os.path.getsize(local_path) > 0:
                                    downloaded_files.append({
                                        'name': file_name,
                                        'path': local_path
                                    })
                                    copied_files.add(file_name)
                                    BuiltIn().log(f"Downloaded: {file_name} from container to {local_path}", level='INFO')
                        
                        if downloaded_files:
                            return downloaded_files
                    else:
                        BuiltIn().log(f"Error checking container files: {result.stderr}", level='WARN')
                else:
                    # Local mode: Check local directory directly
                    pdf_files = [f for f in os.listdir(download_dir) if f.endswith('.pdf')]
                    BuiltIn().log(f"[{elapsed:.1f}s] Checking local directory... found {len(pdf_files)} PDF(s)", level='DEBUG')
                    
                    if pdf_files:
                        for file_name in pdf_files:
                            # Skip if already processed
                            if file_name in copied_files:
                                continue
                            
                            file_path = os.path.join(download_dir, file_name)
                            # Check if file is complete (not being written)
                            if os.path.getsize(file_path) > 0:
                                downloaded_files.append({
                                    'name': file_name,
                                    'path': file_path
                                })
                                copied_files.add(file_name)
                                BuiltIn().log(f"Found downloaded file: {file_name} at {file_path}", level='INFO')
                        
                        if downloaded_files:
                            return downloaded_files
                
                time.sleep(2)
                
            except subprocess.TimeoutExpired:
                BuiltIn().log(f"Timeout checking container files", level='WARN')
                time.sleep(2)
            except Exception as e:
                BuiltIn().log(f"Error checking downloads: {e}", level='WARN')
                import traceback
                BuiltIn().log(traceback.format_exc(), level='DEBUG')
                time.sleep(2)
        
        raise TimeoutError(f"No file downloaded within {timeout} seconds")
