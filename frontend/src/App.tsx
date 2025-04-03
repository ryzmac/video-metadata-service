import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';

interface Video {
  id: string;
  title: string;
  description: string;
  duration: number;
  uploadDate: string;
}

const App: React.FC = () => {
  const [videos, setVideos] = useState<Video[]>([]);
  const [token, setToken] = useState<string>('');

  useEffect(() => {
  const login = async () => {
    try {
      const response = await axios.post('http://localhost:8080/login', {
        username: 'demo',
        password: 'password',
      });
      console.log('Login response:', response.data); // Add this
      setToken(response.data.token);
    } catch (error) {
      console.error('Login failed:', error); // Already there
    }
  };
  login();
}, []);

  useEffect(() => {
    if (token) {
      const fetchVideos = async () => {
        try {
          const response = await axios.get('http://localhost:8080/videos', {
            headers: { Authorization: `Bearer ${token}` },
          });
          setVideos(response.data);
        } catch (error) {
          console.error('Failed to fetch videos:', error);
        }
      };
      fetchVideos();
    }
  }, [token]);

  return (
    <div className="App">
      <h1>Video Metadata</h1>
      {token ? (
        <ul>
          {videos.length > 0 ? (
            videos.map((video) => (
              <li key={video.id}>
                <strong>{video.title}</strong> - {video.duration} seconds
                <br />
                <small>{video.description}</small>
                <br />
                <small>Uploaded: {video.uploadDate}</small>
              </li>
            ))
          ) : (
            <p>No videos available.</p>
          )}
        </ul>
      ) : (
        <p>Logging in...</p>
      )}
    </div>
  );
};

export default App;
