import { FormEvent, useState } from 'react';
import { Link, Navigate, useNavigate } from 'react-router-dom';
import { Alert, Button, TextField, Typography } from '@mui/material';
import { useAuth } from '../context/AuthContext';

export default function RegisterPage() {
  const { register, user } = useAuth();
  const navigate = useNavigate();
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  if (user) return <Navigate to="/" replace />;

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      await register(name, email, password);
      navigate('/');
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Registration failed');
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
              Create account
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Join the operations workspace
            </Typography>
          </div>
        </div>
        {error && (
          <Alert severity="error" className="mb-4">
            {error}
          </Alert>
        )}
        <form className="flex flex-col gap-4" onSubmit={onSubmit}>
          <TextField label="Name" value={name} onChange={(e) => setName(e.target.value)} required fullWidth />
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
            Register
          </Button>
        </form>
        <Typography className="mt-5 text-sm" color="text.secondary">
          Already have an account?{' '}
          <Link to="/login" className="link-accent">
            Sign in
          </Link>
        </Typography>
      </div>
    </div>
  );
}
