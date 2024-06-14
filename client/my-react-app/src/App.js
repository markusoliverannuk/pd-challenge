import axios from 'axios';
import './App.css';
import React, { useState } from 'react';

function App() {
  const [username, setUsername] = useState('');
  const [oldGists, setOldGists] = useState(null);
  const [newGists, setNewGists] = useState(null);
  const [loadingSearch, setLoadingSearch] = useState(false); // State for search loading
  const [loadingTrackedUsers, setLoadingTrackedUsers] = useState(false); // State for tracked users loading
  const [showContainer, setShowContainer] = useState(false);
  const [trackedUsers, setTrackedUsers] = useState(null);

  const handleSearch = async () => {
    setLoadingSearch(true); // Start loading
    setShowContainer(true);
    try {
      const trimmedUsername = username.trim(); // Trim whitespace from username
      const response = await axios.get(`https://api-challenge-v0.techwithmarkus.com/user/${trimmedUsername}`);
      setOldGists(response.data.old_gists);
      setNewGists(response.data.new_gists);
    } catch (error) {
      console.error('Error fetching gists:', error);
    }
    setLoadingSearch(false); // Stop loading
  };

  const handleShowTrackedUsers = async () => {
    setLoadingTrackedUsers(true); // Start loading
    try {
      const response = await axios.get(`https://api-challenge-v0.techwithmarkus.com/trackedusers`);
      setTrackedUsers(response.data); 
    } catch (error) {
      console.error('Error fetching tracked users:', error);
    }
    setLoadingTrackedUsers(false); // Stop loading
  };

  return (
    <div className="container">
      <h1 className="title">2-in-1 Github Gist Scraper + Pipedrive Deal Creator</h1>
      <h2 className='whitetextsubtitle'>Simple User Interface</h2>
      <div className="form">
        <input
          type="text"
          className="input"
          placeholder="Enter GitHub Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <button className="button" onClick={handleSearch} disabled={loadingSearch || loadingTrackedUsers}>
          {loadingSearch ? 'Searching...' : 'Scrape gists & create deals'}
        </button>
        <button className="button" onClick={handleShowTrackedUsers} disabled={loadingSearch || loadingTrackedUsers}>
          {loadingTrackedUsers ? 'Loading...' : 'Show previously scanned users'}
        </button>
      </div>
      {loadingSearch || loadingTrackedUsers ? (
        <p className="loadingtext">Loading... (Please allow up to 10 seconds)</p> // Display loading indicator
      ) : (
        <>
          {showContainer && (
            <div className="gists-container">
              {newGists !== null && (
                <div className="user-info">
                  <h2 className="username">Created deals for the following gists from {username}</h2>
                  <div className="gists-list">
                    {newGists.length === 0 ? (
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
              )}
              {oldGists !== null && (
                <div className="user-info">
                  <h2 className="username">Gists from {username} that you've already seen & created deals for</h2>
                  <div className="gists-list">
                    {oldGists.length === 0 ? (
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
              )}
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
        </>
      )}
      <div className="description">
        <h1 className='whitetext'>So what happens when you execute a request?</h1>
        <p>A much more detailed explanation is on Github @ README</p>
        <p>We start right here at <span className='whitetext'>https://challenge.techwithmarkus.com</span>, where you can enter someone's Github username into the search and expect to see their public gists and automatically create Pipedrive deals for them.</p>
        <p>Now let's say you entered someone's github username into the search and pressed the button to scrape gists and create deals.</p>
        <p>What happens now is that we are sending a request to <span className='whitetext'>https://api-challenge-v2.techwithmarkus.com</span>, providing a valid username to complete the endpoint path, for example</p>
        <p className='whitetext'>https://api-challenge-v2.techwithmarkus.com/user/markusoliverannuk.</p>
        <p>Simply put, after going through the DNS, our request reaches a load balancer and finally the ec2 instances running our application.</p>
        <p>I will now let the following illustrations do the rest of the talking.</p>
        <img src="/schemas/architecture.png" alt="Architecture Diagram" className="media-images" />
        <p>We've managed to make our way through the load balancer and have reached the EC2 instances.</p>
        <img src="/schemas/instances.png" alt="instances" className="media-images" />
        <p>Now let's see what really happens once our request reaches our EC2 instance..</p>
        <img src="/schemas/insidemachine.png" alt="instances" className="media-images" />
        <p>Hopefully these illustrations gave you somewhat of an understanding of how it works. As I said though, a much more detailed explanation is in the README file located inside my Github Repository.</p>
        <p className='whitetextauthor'>- Markus</p>

      </div>
    </div>
  );
}

export default App;
