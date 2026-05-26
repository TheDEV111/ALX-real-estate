import Topbar from "../components/Topbar";
import Navbar from "../components/Navbar";
import PropertyTypeFilter from "../components/PropertyTypeFilter";
import Hero from "../components/Hero";
import PropertyListing from "../components/PropertyListingFilter";
import { PropertyListings } from "../components/PropertyListings";

function Homepage() {
  return (
    <div className="">
      <Topbar />
      <Navbar />
      <PropertyTypeFilter />
      <Hero />
      <PropertyListing />
      <PropertyListings />
    </div>
  );
}

export default Homepage;
