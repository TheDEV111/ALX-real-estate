import Topbar from "../components/Topbar";
import Navbar from "../components/Navbar";
import PropertyTypeFilter from "../components/PropertyTypeFilter";
import Hero from "../components/Hero";

function Homepage() {
  return (
    
      <div className="">
        <Topbar />
        <Navbar />

        <PropertyTypeFilter />
        <Hero />
      </div>
    
  );
}

export default Homepage