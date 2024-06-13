import axios from 'axios';
import './App.css';
import React, { useState, } from 'react';


function App() {
  const [username, setUsername] = useState('');
  const [oldGists, setOldGists] = useState(null);
  const [newGists, setNewGists] = useState(null);
  const [loading, setLoading] = useState(false);
  const [showContainer, setShowContainer] = useState(false);
  const [trackedUsers, setTrackedUsers] = useState(null); 





  const handleSearch = async () => {
    setLoading(true);
    setShowContainer(true);
    try {
      const response = await axios.get(`http://localhost:8050/user/${username}`);
      setOldGists(response.data.old_gists);
      setNewGists(response.data.new_gists);
    } catch (error) {
      console.error('Error fetching gists:', error);
    }
    setLoading(false);
  };

  const handleShowTrackedUsers = async () => {
    try {
      const response = await axios.get(`http://localhost:8050/trackedusers`);
      setTrackedUsers(response.data); 
    } catch (error) {
      console.error('Error fetching tracked users:', error);
    }
  };

  return (
    <div className="container">
      <h1 className="title">2-in-1 Github Gist Scraper + Pipedrive Deal Creator</h1>
      <div className="form">
        <input
          type="text"
          className="input"
          placeholder="Enter GitHub Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <button className="button" onClick={handleSearch} disabled={loading}>
          {loading ? 'Searching...' : 'Scrape gists & create deals'}
        </button>
        <button className="button" onClick={handleShowTrackedUsers}>
          Show previously scanned users
        </button> 
      </div>
      {showContainer && (
        <div className="gists-container">
          <div className="user-info">
            <h2 className="username">Created deals for the following gists from {username}</h2>
            <div className="gists-list">
              {newGists === null ? (
                <p className="loading">No new gists since the last time you visited. Try again later! :)</p>
              ) : newGists.length === 0 ? (
                <p className="no-gists">No new gists</p>
              ) : (
                <ul className="gist-list">
                  {newGists.map((gist) => (
                    <li key={gist.id} className="gist-item">
                      <p className="gist-description">{gist.description}</p>
                      <ul className="files-list">
                        {gist.files.map((file, index) => (
                          <li key={index}>
                            <a href={file} className="file-link" target="_blank" rel="noreferrer">
                              {file}
                            </a>
                          </li>
                        ))}
                      </ul>
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </div>
          <div className="user-info">
            <h2 className="username">Gists from {username} that you've already seen & created deals for</h2>
            <div className="gists-list">
              {oldGists === null ? (
                <p className="loading">This is probably the first time you've searched for {username}'s gists.</p>
              ) : oldGists.length === 0 ? (
                <p className="no-gists">No old gists</p>
              ) : (
                <ul className="gist-list">
                  {oldGists.map((gist) => (
                    <li key={gist.id} className="gist-item">
                      <p className="gist-description">{gist.description}</p>
                      <ul className="files-list">
                        {gist.files.map((file, index) => (
                          <li key={index}>
                            <a href={file} className="file-link" target="_blank" rel="noreferrer">
                              {file}
                            </a>
                          </li>
                        ))}
                      </ul>
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </div>
        </div>
      )}
      {trackedUsers && (
        <div className="tracked-users-container">
          <h2>Previously tracked users:</h2>
          <ul>
            {trackedUsers.map((user, index) => (
              <li key={index}>{user}</li>
            ))}
          </ul>
        </div>
      )} 
      <div className='description'>
        <h1>So what happens when you press 'Search'?</h1>
        <p>A much more detailed explanation is on Github @ README</p>
        <img src="/schemas/architecture.png" alt="Architecture Diagram"  className='media-images'/>
        <img src="/schemas/instances.png" alt="instances"  className='media-images'/>
      </div>
    </div>
  );
}

export default App;
