import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  Box,
  Heading,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Button,
} from "@chakra-ui/react";
import axios from "axios";

interface Car {
  _id: string;
  make: string;
  model: string;
  year: number;
  type: string;
  available: boolean;
}

const CarList = () => {
  const queryClient = useQueryClient();

  const { data: cars, isLoading, error } = useQuery<Car[]>({
    queryKey: ['cars'],
    queryFn: () => axios.get("/api/cars",{ withCredentials: true }).then((res) => res.data),
  });

  const deleteCarMutation = useMutation({
    mutationFn: (id: string) => axios.delete(`/api/cars/${id}`,{ withCredentials: true }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cars'] });
    },
  });

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>An error occurred</div>;

  return (
    <Box width="100%">
      <Heading size="md" mb={4}>
        Car List
      </Heading>
      <Table variant="simple">
        <Thead>
          <Tr>
            <Th>Make</Th>
            <Th>Model</Th>
            <Th>Year</Th>
            <Th>Type</Th>
            <Th>Available</Th>
            <Th>Actions</Th>
          </Tr>
        </Thead>
        <Tbody>
          {cars?.map((car) => (
            <Tr key={car._id}>
              <Td>{car.make}</Td>
              <Td>{car.model}</Td>
              <Td>{car.year}</Td>
              <Td>{car.type}</Td>
              <Td>{car.available ? "Yes" : "No"}</Td>
              <Td>
                <Button
                  colorScheme="red"
                  size="sm"
                  onClick={() => deleteCarMutation.mutate(car._id)}
                >
                  Delete
                </Button>
              </Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
    </Box>
  );
};

export default CarList;