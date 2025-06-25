import React, { useState } from "react";
import "./App.css";

import { OpenFileSelectionDialog, UploadFile } from "../wailsjs/go/main/App";

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

  const handleUpload = async () => {
    UploadFile([selectedFilePath]);
  };

  const handleDownload = async () => {};

  const handleLinkGen = async () => {};

  return (
    <div className="App">
      <header className="App-header">
        <h1>Wails Fileshare Client</h1>

        {error && <p style={{ color: "red" }}>{error}</p>}

        <div className="card">
          <button onClick={handleOpenFile}>Open Single File Dialog</button>
          <p>Selected File: {selectedFilePath}</p>
        </div>

        <div className="card">
          <button onClick={handleUpload}>Upload This File</button>
        </div>

        <div className="card">
          <button onClick={handleLinkGen}>Linkgen This File</button>
        </div>

        <div className="card">
          <button onClick={handleDownload}>Download This File</button>
        </div>
      </header>
    </div>
  );
}

export default App;
