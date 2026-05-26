
import { propertyListings } from "../data/propertyListings";
import { PropertyCard } from "./PropertyCard";

export const PropertyListings: React.FC = () => {
  return (
    <div className=" m-4 md:m-8 py-8">
      <div className=" grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4  md:gap-6 lg:gap-4">
        {propertyListings.map((property) => (
          <PropertyCard key={property.id} property={property} />
        ))}
      </div>
    </div>
  );
};
