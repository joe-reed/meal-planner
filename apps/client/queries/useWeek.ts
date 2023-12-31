import { useQuery } from "react-query";
import { Meal } from "../types/meal";

export async function fetchWeek() {
  const response = await fetch("/api/weeks/this");
  if (!response.ok) {
    throw new Error("Error fetching week");
  }
  return response.json();
}

export default function useWeek() {
  return useQuery<Meal[], Error>("weeks/this", fetchWeek);
}
