import { FormEvent, useState } from 'react';
import { Link } from 'react-router-dom';
import { Alert, Button, Paper, TextField, Typography } from '@mui/material';
import { authApi } from '../services/api';

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState('');
  const [done, setDone] = useState(false);
  const [error, setError] = useState('');

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setError('');
    try {
      await authApi.forgotPassword(email);
      setDone(true);
    } catch {
      setError('Unable to process request');
    }
  }

  return (
    <div className="min-h-screen grid place-items-center p-4">
      <Paper className="glass-panel w-full max-w-md p-8" elevation={0}>
        <Typography variant="h4" className="font-display mb-2">
          Reset password
        </Typography>
        <Typography color="text.secondary" className="mb-6">
          We will email reset instructions if the account exists.
        </Typography>
        {done ? (
          <Alert severity="success">Check your inbox for next steps.</Alert>
        ) : (
          <form className="flex flex-col gap-4" onSubmit={onSubmit}>
            {error && <Alert severity="error">{error}</Alert>}
            <TextField label="Email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} required />
            <Button type="submit" variant="contained">
              Send reset link
            </Button>
          </form>
        )}
        <Link to="/login" className="inline-block mt-4 text-teal-400 hover:underline text-sm">
          Back to login
        </Link>
      </Paper>
    </div>
  );
}
