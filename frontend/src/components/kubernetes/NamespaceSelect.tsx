import { MenuItem, TextField } from '@mui/material';

interface Props {
  value: string;
  options: string[];
  onChange: (ns: string) => void;
}

export default function NamespaceSelect({ value, options, onChange }: Props) {
  return (
    <TextField select size="small" label="Namespace" value={value} onChange={(e) => onChange(e.target.value)}>
      {options.map((ns) => (
        <MenuItem key={ns} value={ns}>
          {ns}
        </MenuItem>
      ))}
    </TextField>
  );
}
