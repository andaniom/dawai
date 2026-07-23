export interface User {
  id: string;
  email: string;
  name: string;
  roles: string[];
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  accessToken: string;
  user: User;
}

export interface ApiResponse<T> {
  success: boolean;
  code: number;
  data: T | null;
  error: {
    message: string;
    type: string;
    details?: Record<string, unknown>;
  } | null;
  meta: {
    timestamp: string;
    path: string;
    version: string;
  } | null;
}

export interface Subject {
  id: string;
  name: string;
  description: string | null;
}

export interface RubricComponent {
  id: string;
  name: string;
  scale_min: number;
  scale_max: number;
  weight: number;
}

export interface Student {
  id: string;
  name: string;
  class: string | null;
}

export interface Assessment {
  id: string;
  student_id: string;
  subject_id: string;
  teacher_id: string;
  feedback: string | null;
  submitted_at: string;
}

export interface AssessmentComponent {
  id: string;
  rubric_component_id: string;
  score: number;
}
