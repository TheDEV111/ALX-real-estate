import Topbar from "../components/Topbar";
import Navbar from "../components/Navbar";
import PropertyTypeFilter from "../components/PropertyTypeFilter";

function Homepage() {
  return (
    <div>
      <Topbar />
      <Navbar />
      
        
        <PropertyTypeFilter />
     
    </div>
  );
}

export default Homepage