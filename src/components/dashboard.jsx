import React, { useEffect, useState } from "react";
import {
  Card,
  CardContent,
  Typography,
  Button,
  CircularProgress,
  Grid,
  Chip,
} from "@mui/material";

const API = process.env.REACT_APP_API;

const Dashboard = () => {
  const [alerts, setAlerts] = useState([]);
  const [loading, setLoading] = useState(false);

  const fetchAlerts = async () => {
    setLoading(true);
    try {
      const res = await fetch(`${API}/alerts`);
      const data = await res.json();
      setAlerts(data.alerts || []);
    } catch (err) {
      console.error("Error fetching alerts", err);
    }
    setLoading(false);
  };

  const remediate = async (alertId) => {
    try {
      await fetch(`${API}/alerts/${alertId}/remediate`, {
        method: "POST",
      });
      fetchAlerts();
    } catch (err) {
      console.error("Remediation failed", err);
    }
  };

  useEffect(() => {
    fetchAlerts();
    const interval = setInterval(fetchAlerts, 5000);
    return () => clearInterval(interval);
  }, []);

  if (loading) return <CircularProgress />;

  return (
    <>
      <Typography variant="h4" gutterBottom>
        SecureCoda Dashboard
      </Typography>

      <Grid container spacing={2}>
        {alerts.map((alert) => (
          <Grid item xs={12} md={6} key={alert.id}>
            <Card>
              <CardContent>
                <Typography variant="h6">{alert.title}</Typography>
                <Typography variant="body2" gutterBottom>
                  {alert.description}
                </Typography>

                <Chip
                  label={alert.type}
                  color={
                    alert.type === "Sensitive Data"
                      ? "error"
                      : alert.type === "Public Share"
                      ? "warning"
                      : "primary"
                  }
                  style={{ marginBottom: "8px" }}
                />

                <Typography variant="caption" display="block">
                  Document: {alert.documentName}
                </Typography>

                <Button
                  variant="contained"
                  color="primary"
                  style={{ marginTop: "1rem" }}
                  onClick={() => remediate(alert.id)}
                >
                  Remediate
                </Button>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      {alerts.length === 0 && (
        <Typography>No alerts found — system is clean 🎉</Typography>
      )}
    </>
  );
};

export default Dashboard;
