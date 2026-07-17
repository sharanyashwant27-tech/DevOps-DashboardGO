import { Button, Stack } from '@mui/material';

interface Props {
  onStart: () => void;
  onStop: () => void;
  onRestart: () => void;
  onLogs: () => void;
  onDelete: () => void;
}

export default function ContainerActions({ onStart, onStop, onRestart, onLogs, onDelete }: Props) {
  return (
    <Stack direction="row" spacing={1} justifyContent="flex-end">
      <Button size="small" onClick={onStart}>Start</Button>
      <Button size="small" onClick={onStop}>Stop</Button>
      <Button size="small" onClick={onRestart}>Restart</Button>
      <Button size="small" onClick={onLogs}>Logs</Button>
      <Button size="small" color="error" onClick={onDelete}>Delete</Button>
    </Stack>
  );
}
