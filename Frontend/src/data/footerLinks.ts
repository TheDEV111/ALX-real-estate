export interface FooterLink {
  label: string;
  href: string;
}

export interface FooterSection {
  title: string;
  links: FooterLink[];
}
   
export const footerSections: FooterSection[] = [
  {
    title: "Explore",
    links: [
      { label: "Apartments in Dubai", href: "#" },
      { label: "Hotels in New York", href: "#" },
      { label: "Villa in Spain", href: "#" },
      { label: "Mansion in Indonesia", href: "#" },
    ],
  },
  {
    title: "Company",
    links: [
      { label: "About us", href: "#" },
      { label: "Blog", href: "#" },
      { label: "Career", href: "#" },
      { label: "Customers", href: "#" },
      { label: "Brand", href: "#" },
    ],
  },
  {
    title: "Help",
    links: [
      { label: "Support", href: "#" },
      { label: "Cancel booking", href: "#" },
      { label: "Refunds Process", href: "#" },
    ],
  },
];

export const footerDescription =
  "ALX is a platform where travelers can discover and book unique, comfortable, and affordable lodging options worldwide. From cozy city apartments and tranquil countryside retreats to exotic beachside villas, ALX connects you with the perfect place to stay for any trip.";

export const footerNote =
  "Some hotel requires you to cancel more than 24 hours before check-in. Details";

export const footerPolicies = [
  { label: "Terms of Service", href: "#" },
  { label: "Policy service", href: "#" },
  { label: "Cookies Policy", href: "#" },
];
