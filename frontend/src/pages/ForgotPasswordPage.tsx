import { FormEvent, useState } from 'react';
import { Link } from 'react-router-dom';
import { Alert, Button, TextField, Typography } from '@mui/material';
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
    <div className="auth-shell">
      <div className="auth-card">
        <div className="flex items-center gap-3 mb-5">
          <span className="brand-mark">DC</span>
          <div>
            <Typography variant="h5" className="font-display">
              Reset password
            </Typography>
            <Typography variant="body2" color="text.secondary">
              We will email reset instructions if the account exists.
            </Typography>
          </div>
        </div>
        {done ? (
          <Alert severity="success">Check your inbox for next steps.</Alert>
        ) : (
          <form className="flex flex-col gap-4" onSubmit={onSubmit}>
            {error && <Alert severity="error">{error}</Alert>}
            <TextField label="Email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} required />
            <Button type="submit" variant="contained" size="large">
              Send reset link
            </Button>
          </form>
        )}
        <Link to="/login" className="link-accent inline-block mt-5 text-sm">
          Back to login
        </Link>
      </div>
    </div>
  );
}
