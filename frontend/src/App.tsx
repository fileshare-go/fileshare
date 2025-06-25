import React, { useState } from "react";
import "./App.css";

import {
  OpenFileSelectionDialog,
  OpenMultipleFilesSelectionDialog,
} from "../wailsjs/go/main/App";

function App() {
  const [selectedFilePath, setSelectedFilePath] = useState<string>("");
  const [selectedFilePaths, setSelectedFilePaths] = useState<string[]>([]);
  const [error, setError] = useState<string>("");

  const handleOpenFile = async () => {
    setError("");
    try {
      const path = await OpenFileSelectionDialog();
      if (path) {
        setSelectedFilePath(path);
      } else {
        setSelectedFilePath("No file selected.");
      }
    } catch (err: any) {
      setError(`Error opening file dialog: ${err.message || String(err)}`);
      setSelectedFilePath("Error!");
    }
  };

  const handleOpenMultipleFiles = async () => {
    setError("");
    try {
      const paths = await OpenMultipleFilesSelectionDialog();
      if (paths && paths.length > 0) {
        setSelectedFilePaths(paths);
      } else {
        setSelectedFilePaths(["No files selected."]);
      }
    } catch (err: any) {
      setError(
        `Error opening multiple files dialog: ${err.message || String(err)}`,
      );
      setSelectedFilePaths(["Error!"]);
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>Wails File Dialog Demo</h1>

        {error && <p style={{ color: "red" }}>{error}</p>}

        <div className="card">
          <button onClick={handleOpenFile}>Open Single File Dialog</button>
          <p>Selected File: {selectedFilePath}</p>
        </div>

        <div className="card">
          <button onClick={handleOpenMultipleFiles}>
            Open Multiple Files Dialog
          </button>
          <p>Selected Files: {selectedFilePaths.join(", ")}</p>
        </div>
      </header>
    </div>
  );
}

export default App;
