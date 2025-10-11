'use client';

import { ProjectStatusColumn } from '@/lib/schemas/projects';
import { Label } from '../ui/label';
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '../ui/select';
import { useTaskForm } from './task-form-context';

type SelectColumnProps = {
  columns: ProjectStatusColumn[];
};

export const SelectColumn = ({ columns }: SelectColumnProps) => {
  const { state } = useTaskForm();
  return (
    <div>
      <Label className="mb-2">Status</Label>
      <Select defaultValue={state.column_id} name="column_id">
        <SelectTrigger className="w-[180px] cursor-pointer">
          <SelectValue placeholder="Select status" />
        </SelectTrigger>
        <SelectContent>
          <SelectGroup>
            {columns.map(({ id, name }) => (
              <SelectItem
                key={id}
                value={id.toString()}
                className="cursor-pointer"
              >
                {name}
              </SelectItem>
            ))}
          </SelectGroup>
        </SelectContent>
      </Select>
    </div>
  );
};
