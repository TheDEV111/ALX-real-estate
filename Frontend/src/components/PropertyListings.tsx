
import { propertyListings } from "../data/propertyListings";
import { Button } from "./Button";
import { PropertyCard } from "./PropertyCard";

export const PropertyListings: React.FC = () => {
  return (
    <div className=" m-4 md:m-8 py-8">
      <div className=" grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4  md:gap-6 lg:gap-4">
        {propertyListings.map((property) => (
          <PropertyCard key={property.id} property={property} />
        ))}
      </div>
      <div className="flex flex-col justify-center items-center p-[50px] font-[Quicksand] gap-[10px]  mt-[70px]">
        <Button className="bg-[#161117] text-white px-[32px] py-[13px] font-medium text-[20px]">
          Show more
        </Button>
        <p className="font-medium text-[20px]">Click to see more listings</p>
      </div>
      
    </div>
  );
};
