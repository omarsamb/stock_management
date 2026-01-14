import { Doughnut } from 'react-chartjs-2';
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';

ChartJS.register(ArcElement, Tooltip, Legend);

export const StockPieChart = ({ data }: { data: any[] }) => {
  if (!data || data.length === 0) return <p>Aucune donnée de catégorie disponible.</p>;

  const chartData = {
    labels: data.map(d => d.category_name),
    datasets: [
      {
        data: data.map(d => d.total_value),
        backgroundColor: [
          '#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF', '#FF9F40', '#4BC0C0', '#FF9F40'
        ],
        hoverBackgroundColor: [
          '#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF', '#FF9F40', '#4BC0C0', '#FF9F40'
        ],
      },
    ],
  };
  
  const options = {
    responsive: true,
    plugins: {
      legend: {
        position: 'bottom' as const,
      },
    },
  };

  return <Doughnut data={chartData} options={options} />;
};
