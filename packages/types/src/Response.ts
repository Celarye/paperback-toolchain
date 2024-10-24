import type { Cookie } from "./Cookie"

export interface Response {
  readonly url: string;
  readonly headers: Record<string, string>;
  readonly status: number;
  readonly mimeType?: string;

  /// This is only the new cookies set via the Set-Cookie header
  readonly cookies: Cookie[];
}