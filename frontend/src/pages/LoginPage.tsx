import { FormEvent, useState } from 'react';
import { Link, Navigate, useNavigate } from 'react-router-dom';
import { Alert, Button, TextField, Typography } from '@mui/material';
import { useAuth } from '../context/AuthContext';

export default function LoginPage() {
  const { login, user } = useAuth();
  const navigate = useNavigate();
  const [email, setEmail] = useState('admin@devops.local');
  const [password, setPassword] = useState('Admin@12345');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  if (user) return <Navigate to="/" replace />;

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      await login(email, password);
      navigate('/');
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Login failed');
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="auth-shell">
      <div className="auth-card">
        <div className="flex items-center gap-3 mb-5">
          <span className="brand-mark">DC</span>
          <div>
            <Typography variant="h5" className="font-display">
              DevOps Command Center
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Sign in to your colorful ops workspace
            </Typography>
          </div>
        </div>
        {error && (
          <Alert severity="error" className="mb-4">
            {error}
          </Alert>
        )}
        <form className="flex flex-col gap-4" onSubmit={onSubmit}>
          <TextField label="Email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} required fullWidth />
          <TextField
            label="Password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            fullWidth
          />
          <Button type="submit" variant="contained" size="large" disabled={loading}>
            {loading ? 'Signing in...' : 'Sign in'}
          </Button>
        </form>
        <div className="mt-5 flex justify-between text-sm">
          <Link to="/register" className="link-accent">
            Create account
          </Link>
          <Link to="/forgot-password" className="link-accent">
            Forgot password
          </Link>
        </div>
      </div>
    </div>
  );
}
