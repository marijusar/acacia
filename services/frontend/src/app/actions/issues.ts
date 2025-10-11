'use server';

import { createIssueRequestBody } from '@/lib/schemas/issues';
import { Issue, updateIssueRequestBody } from '@/lib/schemas/projects';
import { issuesService } from '@/lib/services/issues-service';
import { revalidatePath } from 'next/cache';

export async function createIssueAction(formData: FormData) {
  const fields = Object.fromEntries(formData);

  await issuesService.createIssue(createIssueRequestBody.parse(fields));

  revalidatePath('/projects', 'layout');
}

export async function updateIssue(issue: Issue) {
  await issuesService.updateIssue(issue);

  revalidatePath('/projects', 'layout');
}

export async function updateIssueForm(formData: FormData) {
  const entries = Object.fromEntries(formData);
  await issuesService.updateIssue(updateIssueRequestBody.parse(entries));

  revalidatePath('/projects', 'layout');
}
