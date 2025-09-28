'use server';

import { createIssueRequestBody } from '@/lib/schemas/issues';
import { issuesService } from '@/lib/services/issues-service';
import { revalidatePath } from 'next/cache';

export async function createIssueAction(formData: FormData) {
  console.log(formData);
  const fields = Object.fromEntries(formData);
  console.log({ fields });

  await issuesService.createIssue(createIssueRequestBody.parse(fields));

  revalidatePath('/projects', 'layout');
}
