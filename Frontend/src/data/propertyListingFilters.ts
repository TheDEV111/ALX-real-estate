export interface FilterOption {
  id: number;
  name: string;
  icon?: string;
}

export const filterOptions: FilterOption[] = [
  { id: 1, name: "All" },
  { id: 2, name: "Top Villa" },
  { id: 3, name: "Free Reschedule" },
  { id: 4, name: "Book Now, Pay later" },
  { id: 5, name: "Self CheckIn" },
  { id: 6, name: "Instant Book" },
];

export const sortOptions = [
  "Highest Price",
  "Lowest Price",
  "Most Popular",
  "Newest",
];
