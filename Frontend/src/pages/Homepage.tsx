import Topbar from "../components/Topbar";
import Navbar from "../components/Navbar";
import PropertyTypeFilter from "../components/PropertyTypeFilter";
import Hero from "../components/Hero";
import PropertyListing from "../components/PropertyListingFilter";
import { PropertyListings } from "../components/PropertyListings";
import Footer from "../components/Footer";

function Homepage() {
  return (
    <div className="">
      <Topbar />
      <Navbar />
      <PropertyTypeFilter />
      <Hero />
      <PropertyListing />
      <PropertyListings />
      <div className="bg-[#34967C] h-6"></div>
      <Footer />
    </div>
  );
}

export default Homepage;
