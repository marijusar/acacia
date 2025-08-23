import { returnsTrue } from '#acacia/utils/example.ts';
import { describe, it, expect } from 'vitest';

describe('example', () => {
  it('should work', () => {
    expect(returnsTrue()).toEqual(true);
  });
});
