import { useState } from 'react';
import { 
  Container, 
  Typography, 
  Paper, 
  Box,
  TextField,
  Button
} from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer
} from 'recharts';

interface DORAMetrics {
  DeploymentFrequency: number;
  LeadTimeForChanges: number;
  ChangeFailureRate: number;
}

function App() {
  const [startDate, setStartDate] = useState<Date | null>(new Date(Date.now() - 30 * 24 * 60 * 60 * 1000));
  const [endDate, setEndDate] = useState<Date | null>(new Date());
  const [repoPath, setRepoPath] = useState<string>('');
  const [metrics, setMetrics] = useState<DORAMetrics | null>(null);

  const fetchMetrics = async () => {
    if (!startDate || !endDate || !repoPath) return;

    try {
      const response = await fetch('http://localhost:8099/api/metrics', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          startDate: startDate.toISOString().split('T')[0],
          endDate: endDate.toISOString().split('T')[0],
          repoPath: repoPath,
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to fetch metrics');
      }

      const data = await response.json();
      setMetrics(data);
    } catch (error) {
      console.error('Error fetching metrics:', error);
    }
  };

  console.log('Current metrics state:', metrics);
  console.log('Type of metrics state:', typeof metrics);
  console.log('Is object:', typeof metrics === 'object');

  const chartData = metrics ? [
    {
      name: 'Deployment Frequency',
      value: metrics.DeploymentFrequency || 0,
    },
    {
      name: 'Lead Time (hours)',
      value: metrics.LeadTimeForChanges || 0,
    },
    {
      name: 'Change Failure Rate (%)',
      value: metrics.ChangeFailureRate || 0,
    },
  ] : [];

  console.log('Chart data:', chartData);

  return (
    <LocalizationProvider dateAdapter={AdapterDateFns}>
      <Container maxWidth="lg">
        <Box sx={{ my: 4 }}>
          <Typography variant="h4" component="h1" gutterBottom>
            DORA Metrics Dashboard
          </Typography>

          <Paper sx={{ p: 3, mb: 3 }}>
            <Box sx={{ display: 'flex', gap: 2, flexDirection: 'column' }}>
              <Box sx={{ width: '100%' }}>
                <TextField
                  fullWidth
                  label="Repository Path"
                  value={repoPath}
                  onChange={(e) => setRepoPath(e.target.value)}
                  placeholder="Enter the path to your git repository"
                />
              </Box>
              <Box sx={{ 
                display: 'flex', 
                gap: 2,
                flexDirection: { xs: 'column', sm: 'row' },
                alignItems: { xs: 'stretch', sm: 'flex-start' }
              }}>
                <Box sx={{ flex: 1 }}>
                  <DatePicker
                    label="Start Date"
                    value={startDate}
                    onChange={(newValue) => setStartDate(newValue)}
                    sx={{ width: '100%' }}
                  />
                </Box>
                <Box sx={{ flex: 1 }}>
                  <DatePicker
                    label="End Date"
                    value={endDate}
                    onChange={(newValue) => setEndDate(newValue)}
                    sx={{ width: '100%' }}
                  />
                </Box>
                <Box sx={{ 
                  flex: 1,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center'
                }}>
                  <Button 
                    variant="contained" 
                    onClick={fetchMetrics} 
                    disabled={!repoPath}
                    sx={{ width: '100%' }}
                  >
                    Fetch Metrics
                  </Button>
                </Box>
              </Box>
            </Box>
          </Paper>

          {metrics && (
            <Paper sx={{ p: 3 }}>
              <Typography variant="h6" gutterBottom>
                DORA Metrics Overview
              </Typography>
              <Box sx={{ height: 400 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={chartData}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="name" />
                    <YAxis />
                    <Tooltip />
                    <Legend />
                    <Bar dataKey="value" fill="#8884d8" />
                  </BarChart>
                </ResponsiveContainer>
              </Box>
            </Paper>
          )}
        </Box>
      </Container>
    </LocalizationProvider>
  );
}

export default App;
